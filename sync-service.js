// Background Sync Service for Supabase
// Ensures all signals are synced and enriched with data

class SyncService {
    constructor() {
        this.syncInterval = 30000; // 30 seconds
        this.isRunning = false;
        this.lastSync = null;
    }

    // Start background sync
    start() {
        if (this.isRunning) {
            console.log('🔄 Sync service already running');
            return;
        }

        this.isRunning = true;
        console.log('🚀 Starting background sync service...');
        
        // Initial sync
        this.performSync();
        
        // Periodic sync
        this.intervalId = setInterval(() => {
            this.performSync();
        }, this.syncInterval);
    }

    // Stop background sync
    stop() {
        if (this.intervalId) {
            clearInterval(this.intervalId);
            this.isRunning = false;
            console.log('⏹️ Sync service stopped');
        }
    }

    // Perform full sync
    async performSync() {
        try {
            console.log('🔄 Starting sync...');
            
            // Step 1: Sync missing signals from localStorage
            await this.syncMissingSignals();
            
            // Step 2: Enrich signals with missing data
            await this.enrichSignals();
            
            // Step 3: Update live prices for pending signals
            await this.updateLivePrices();
            
            this.lastSync = new Date();
            console.log('✅ Sync completed at', this.lastSync.toLocaleTimeString());
            
        } catch (error) {
            console.error('❌ Sync error:', error);
        }
    }

    // Sync missing signals from localStorage to Supabase
    async syncMissingSignals() {
        if (typeof SupabaseDB === 'undefined') {
            console.warn('⚠️ Supabase not available');
            return;
        }

        try {
            // Get all signals from localStorage
            const localSignals = JSON.parse(localStorage.getItem('tradingSignals') || '[]');
            
            if (localSignals.length === 0) {
                console.log('📭 No local signals to sync');
                return;
            }

            // Get all signals from Supabase
            const { data: supabaseSignals, error } = await supabase
                .from('trading_signals')
                .select('signal_id');

            if (error) throw error;

            const supabaseIds = new Set(supabaseSignals.map(s => s.signal_id));
            
            // Find missing signals
            const missingSignals = localSignals.filter(s => 
                !supabaseIds.has(s.id?.toString())
            );

            if (missingSignals.length === 0) {
                console.log('✅ All signals synced');
                return;
            }

            console.log(`📤 Syncing ${missingSignals.length} missing signals...`);

            // Sync each missing signal
            for (const signal of missingSignals) {
                await this.saveSignalToSupabase(signal);
            }

            console.log(`✅ Synced ${missingSignals.length} signals`);

        } catch (error) {
            console.error('❌ Error syncing signals:', error);
        }
    }

    // Save signal to Supabase with enriched data
    async saveSignalToSupabase(signal) {
        try {
            const enrichedSignal = {
                signal_id: signal.id?.toString() || Date.now().toString(),
                signal_type: signal.type || 'BUY',
                symbol: signal.symbol || 'BTCUSDT',
                entry_price: signal.entry || 0,
                stop_loss: signal.stopLoss || 0,
                tp1: signal.tp1 || 0,
                tp2: signal.tp2 || 0,
                tp3: signal.tp3 || 0,
                strength: signal.strength || 50,
                pattern_type: signal.patternType || this.detectPattern(signal),
                pattern_confidence: signal.patternConfidence || signal.strength || 50,
                kill_zone: signal.killZone || this.detectKillZone(),
                session_type: signal.sessionType || this.detectSession(),
                status: signal.status || 'pending',
                trailing_stop_active: signal.trailingStopActive || false,
                exit_price: signal.exitPrice || null,
                exit_reason: signal.exitReason || null,
                profit_percent: signal.profitPercent || null,
                live_price: signal.livePrice || null
            };

            const { error } = await supabase
                .from('trading_signals')
                .insert([enrichedSignal]);

            if (error) {
                // If duplicate, update instead
                if (error.code === '23505') {
                    await this.updateSignalInSupabase(enrichedSignal);
                } else {
                    throw error;
                }
            }

        } catch (error) {
            console.error('❌ Error saving signal:', error);
        }
    }

    // Update existing signal in Supabase
    async updateSignalInSupabase(signal) {
        try {
            const { error } = await supabase
                .from('trading_signals')
                .update({
                    pattern_type: signal.pattern_type,
                    pattern_confidence: signal.pattern_confidence,
                    kill_zone: signal.kill_zone,
                    session_type: signal.session_type,
                    status: signal.status,
                    exit_price: signal.exit_price,
                    exit_reason: signal.exit_reason,
                    profit_percent: signal.profit_percent,
                    live_price: signal.live_price
                })
                .eq('signal_id', signal.signal_id);

            if (error) throw error;

        } catch (error) {
            console.error('❌ Error updating signal:', error);
        }
    }

    // Enrich signals with missing data
    async enrichSignals() {
        if (typeof SupabaseDB === 'undefined') return;

        try {
            // Get signals with null values
            const { data: signals, error } = await supabase
                .from('trading_signals')
                .select('*')
                .or('pattern_type.is.null,kill_zone.is.null,session_type.is.null');

            if (error) throw error;
            if (!signals || signals.length === 0) return;

            console.log(`🔧 Enriching ${signals.length} signals...`);

            for (const signal of signals) {
                const updates = {};

                // Fill pattern_type
                if (!signal.pattern_type) {
                    updates.pattern_type = this.detectPattern(signal);
                }

                // Fill pattern_confidence
                if (!signal.pattern_confidence) {
                    updates.pattern_confidence = signal.strength || 50;
                }

                // Fill kill_zone
                if (!signal.kill_zone) {
                    updates.kill_zone = this.detectKillZone(new Date(signal.created_at));
                }

                // Fill session_type
                if (!signal.session_type) {
                    updates.session_type = this.detectSession(new Date(signal.created_at));
                }

                // Update if we have changes
                if (Object.keys(updates).length > 0) {
                    await supabase
                        .from('trading_signals')
                        .update(updates)
                        .eq('signal_id', signal.signal_id);
                }
            }

            console.log(`✅ Enriched ${signals.length} signals`);

        } catch (error) {
            console.error('❌ Error enriching signals:', error);
        }
    }

    // Update live prices for pending signals
    async updateLivePrices() {
        if (typeof SupabaseDB === 'undefined') return;

        try {
            // Get pending signals
            const { data: signals, error } = await supabase
                .from('trading_signals')
                .select('*')
                .eq('status', 'pending');

            if (error) throw error;
            if (!signals || signals.length === 0) return;

            // Get current price from Binance
            const response = await fetch('https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT');
            const data = await response.json();
            const currentPrice = parseFloat(data.price);

            // Update each signal
            for (const signal of signals) {
                await supabase
                    .from('trading_signals')
                    .update({ live_price: currentPrice })
                    .eq('signal_id', signal.signal_id);
            }

        } catch (error) {
            // Silently fail for price updates
        }
    }

    // Detect pattern based on signal characteristics
    detectPattern(signal) {
        if (!signal) return 'Trend Following';

        const patterns = [
            'Engulfing',
            'Pin Bar',
            'Inside Bar',
            'Trend Following',
            'Breakout',
            'Reversal',
            'Continuation'
        ];

        // Simple pattern detection based on strength
        const strength = signal.strength || 50;
        
        if (strength >= 80) return 'Engulfing';
        if (strength >= 70) return 'Pin Bar';
        if (strength >= 60) return 'Breakout';
        if (strength >= 50) return 'Trend Following';
        return 'Continuation';
    }

    // Detect kill zone based on time
    detectKillZone(date = new Date()) {
        const hour = date.getUTCHours();
        
        // London Open: 8:00-10:00 UTC
        if (hour >= 8 && hour < 10) return 'London Open';
        
        // London Close: 16:00-18:00 UTC
        if (hour >= 16 && hour < 18) return 'London Close';
        
        // New York Open: 13:00-15:00 UTC
        if (hour >= 13 && hour < 15) return 'New York Open';
        
        // New York Close: 21:00-23:00 UTC
        if (hour >= 21 && hour < 23) return 'New York Close';
        
        // Asian Session: 0:00-8:00 UTC
        if (hour >= 0 && hour < 8) return 'Asian Session';
        
        return 'Off Hours';
    }

    // Detect trading session
    detectSession(date = new Date()) {
        const hour = date.getUTCHours();
        
        if (hour >= 0 && hour < 8) return 'Asian';
        if (hour >= 8 && hour < 16) return 'London';
        if (hour >= 13 && hour < 21) return 'New York';
        
        return 'Off Hours';
    }

    // Get sync status
    getStatus() {
        return {
            isRunning: this.isRunning,
            lastSync: this.lastSync,
            syncInterval: this.syncInterval
        };
    }
}

// Create global instance
window.syncService = new SyncService();

// Auto-start on page load
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        setTimeout(() => window.syncService.start(), 2000);
    });
} else {
    setTimeout(() => window.syncService.start(), 2000);
}

console.log('📦 Sync Service loaded. Use window.syncService to control.');
