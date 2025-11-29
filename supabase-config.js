// Supabase Configuration
const SUPABASE_URL = 'https://xlxugbqxfrrwutxecwug.supabase.co';
const SUPABASE_ANON_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InhseHVnYnF4ZnJyd3V0eGVjd3VnIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ0MTM0MDQsImV4cCI6MjA3OTk4OTQwNH0.LAQKcsWOOTOeAcKh9zRF9l740aDp0Ki1cws8oG3yQZU';

// Initialize Supabase client
const supabase = window.supabase.createClient(SUPABASE_URL, SUPABASE_ANON_KEY);

// Database operations
const SupabaseDB = {
    // Save a new signal
    async saveSignal(signal) {
        try {
            const { data, error } = await supabase
                .from('trading_signals')
                .insert([{
                    signal_id: signal.id.toString(),
                    signal_type: signal.type,
                    symbol: signal.symbol,
                    entry_price: signal.entry,
                    stop_loss: signal.stopLoss,
                    tp1: signal.tp1,
                    tp2: signal.tp2,
                    tp3: signal.tp3,
                    strength: signal.strength,
                    pattern_type: signal.patternType || null,
                    pattern_confidence: signal.patternConfidence || null,
                    kill_zone: signal.killZone || null,
                    session_type: signal.sessionType || null,
                    status: signal.status || 'pending',
                    trailing_stop_active: signal.trailingStopActive || false
                }])
                .select();

            if (error) throw error;
            console.log('✅ Signal saved to Supabase:', data);
            return data[0];
        } catch (error) {
            console.error('❌ Error saving signal:', error);
            throw error;
        }
    },

    // Update signal status (win/loss)
    async updateSignalStatus(signalId, updates) {
        try {
            const { data, error } = await supabase
                .from('trading_signals')
                .update({
                    status: updates.status,
                    exit_price: updates.exitPrice,
                    exit_reason: updates.exitReason,
                    exit_time: new Date().toISOString(),
                    profit_percent: updates.profitPercent,
                    profit_pips: updates.profitPips,
                    holding_time_minutes: updates.holdingTimeMinutes,
                    trailing_stop_price: updates.trailingStopPrice
                })
                .eq('signal_id', signalId.toString())
                .select();

            if (error) throw error;
            console.log('✅ Signal updated:', data);
            return data[0];
        } catch (error) {
            console.error('❌ Error updating signal:', error);
            throw error;
        }
    },

    // Update live price for pending signals
    async updateLivePrice(signalId, livePrice, trailingStop = null) {
        try {
            const updates = {
                live_price: livePrice
            };
            
            if (trailingStop) {
                updates.trailing_stop_price = trailingStop;
                updates.trailing_stop_active = true;
            }

            const { data, error } = await supabase
                .from('trading_signals')
                .update(updates)
                .eq('signal_id', signalId.toString())
                .select();

            if (error) throw error;
            return data[0];
        } catch (error) {
            console.error('❌ Error updating live price:', error);
            throw error;
        }
    },

    // Get all signals
    async getAllSignals() {
        try {
            const { data, error } = await supabase
                .from('trading_signals')
                .select('*')
                .order('created_at', { ascending: false });

            if (error) throw error;
            return data;
        } catch (error) {
            console.error('❌ Error fetching signals:', error);
            throw error;
        }
    },

    // Get pending signals
    async getPendingSignals() {
        try {
            const { data, error } = await supabase
                .from('trading_signals')
                .select('*')
                .eq('status', 'pending')
                .order('created_at', { ascending: false });

            if (error) throw error;
            return data;
        } catch (error) {
            console.error('❌ Error fetching pending signals:', error);
            throw error;
        }
    },

    // Get analytics
    async getAnalytics() {
        try {
            const { data, error } = await supabase
                .from('signal_analytics')
                .select('*');

            if (error) throw error;
            return data;
        } catch (error) {
            console.error('❌ Error fetching analytics:', error);
            throw error;
        }
    },

    // Delete a signal
    async deleteSignal(signalId) {
        try {
            const { error } = await supabase
                .from('trading_signals')
                .delete()
                .eq('signal_id', signalId.toString());

            if (error) throw error;
            console.log('✅ Signal deleted');
            return true;
        } catch (error) {
            console.error('❌ Error deleting signal:', error);
            throw error;
        }
    },

    // Get signals by filter
    async getSignalsByFilter(filters) {
        try {
            let query = supabase
                .from('trading_signals')
                .select('*');

            if (filters.status) {
                query = query.eq('status', filters.status);
            }
            if (filters.signalType) {
                query = query.eq('signal_type', filters.signalType);
            }
            if (filters.killZone) {
                query = query.eq('kill_zone', filters.killZone);
            }
            if (filters.symbol) {
                query = query.eq('symbol', filters.symbol);
            }
            if (filters.dateFrom) {
                query = query.gte('created_at', filters.dateFrom);
            }
            if (filters.dateTo) {
                query = query.lte('created_at', filters.dateTo);
            }

            query = query.order('created_at', { ascending: false });

            const { data, error } = await query;
            if (error) throw error;
            return data;
        } catch (error) {
            console.error('❌ Error filtering signals:', error);
            throw error;
        }
    },

    // Subscribe to real-time updates
    subscribeToSignals(callback) {
        const subscription = supabase
            .channel('trading_signals_changes')
            .on('postgres_changes', 
                { event: '*', schema: 'public', table: 'trading_signals' },
                (payload) => {
                    console.log('🔄 Real-time update:', payload);
                    callback(payload);
                }
            )
            .subscribe();

        return subscription;
    }
};

// Export for use in other files
window.SupabaseDB = SupabaseDB;
