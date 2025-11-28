// WhatsApp Alert Configuration
// Follow WHATSAPP_SETUP.md for instructions

// ═══════════════════════════════════════════════════════════════════
//  CONFIGURATION
// ═══════════════════════════════════════════════════════════════════

const WHATSAPP_CONFIG = {
    enabled: false,  // Set to true after setup
    
    // CallMeBot Configuration (FREE)
    phone: '+1234567890',  // Your phone with country code (e.g., +1 for US, +91 for India)
    apiKey: 'YOUR_API_KEY_HERE',  // Get from CallMeBot (see WHATSAPP_SETUP.md)
    
    // Message Settings
    includeChart: false,  // Future: Include chart image
    soundEnabled: true,   // Play sound with WhatsApp alert
};

// ═══════════════════════════════════════════════════════════════════
//  WHATSAPP ALERT FUNCTION
// ═══════════════════════════════════════════════════════════════════

async function sendWhatsAppAlert(signal) {
    if (!WHATSAPP_CONFIG.enabled || !signal) {
        console.log('WhatsApp alerts disabled');
        return;
    }
    
    if (!WHATSAPP_CONFIG.phone || !WHATSAPP_CONFIG.apiKey) {
        console.error('❌ WhatsApp not configured. Check whatsapp-config.js');
        return;
    }
    
    try {
        // Format message
        const emoji = signal.type === 'BUY' ? '📈' : '📉';
        const message = `${emoji} *${signal.type} SIGNAL*\n\n` +
                       `*Symbol:* ${window.currentSymbol || 'BTCUSDT'}\n` +
                       `*Timeframe:* ${window.currentInterval || '15m'}\n\n` +
                       `*Entry:* ${signal.entry.toFixed(2)}\n` +
                       `*Stop Loss:* ${signal.stopLoss.toFixed(2)}\n\n` +
                       `*Take Profits:*\n` +
                       `TP1: ${signal.targets[0].price.toFixed(2)} (${signal.targets[0].rr.toFixed(1)}R)\n` +
                       `TP2: ${signal.targets[1].price.toFixed(2)} (${signal.targets[1].rr.toFixed(1)}R)\n` +
                       `TP3: ${signal.targets[2].price.toFixed(2)} (${signal.targets[2].rr.toFixed(1)}R)\n\n` +
                       `*Strength:* ${signal.strength}%\n` +
                       `*Risk:* ${signal.riskPercent.toFixed(2)}%\n` +
                       `*Best RR:* ${signal.bestRR.toFixed(2)}:1`;
        
        // Encode message for URL
        const encodedMessage = encodeURIComponent(message);
        
        // CallMeBot API endpoint
        const url = `https://api.callmebot.com/whatsapp.php?phone=${WHATSAPP_CONFIG.phone}&text=${encodedMessage}&apikey=${WHATSAPP_CONFIG.apiKey}`;
        
        console.log('📱 Sending WhatsApp alert...');
        
        // Send request
        const response = await fetch(url, {
            method: 'GET',
            mode: 'no-cors'  // CallMeBot doesn't support CORS
        });
        
        console.log('✅ WhatsApp alert sent!');
        
        // Update status
        const statusEl = document.getElementById('status-message');
        if (statusEl) {
            const originalText = statusEl.textContent;
            statusEl.textContent = '📱 WhatsApp alert sent!';
            statusEl.style.color = '#26a69a';
            setTimeout(() => {
                statusEl.textContent = originalText;
                statusEl.style.color = '#2962ff';
            }, 3000);
        }
        
        return true;
    } catch (error) {
        console.error('❌ WhatsApp send failed:', error);
        return false;
    }
}

// ═══════════════════════════════════════════════════════════════════
//  TELEGRAM ALTERNATIVE (More Reliable)
// ═══════════════════════════════════════════════════════════════════

const TELEGRAM_CONFIG = {
    enabled: false,  // Set to true to use Telegram instead
    botToken: 'YOUR_BOT_TOKEN',  // From @BotFather
    chatId: 'YOUR_CHAT_ID',  // From @userinfobot
};

async function sendTelegramAlert(signal) {
    if (!TELEGRAM_CONFIG.enabled || !signal) return;
    
    try {
        const emoji = signal.type === 'BUY' ? '📈' : '📉';
        const message = `${emoji} *${signal.type} SIGNAL*\n\n` +
                       `*Symbol:* ${window.currentSymbol || 'BTCUSDT'}\n` +
                       `*Timeframe:* ${window.currentInterval || '15m'}\n\n` +
                       `*Entry:* \`${signal.entry.toFixed(2)}\`\n` +
                       `*Stop Loss:* \`${signal.stopLoss.toFixed(2)}\`\n\n` +
                       `*Take Profits:*\n` +
                       `TP1: \`${signal.targets[0].price.toFixed(2)}\` (${signal.targets[0].rr.toFixed(1)}R)\n` +
                       `TP2: \`${signal.targets[1].price.toFixed(2)}\` (${signal.targets[1].rr.toFixed(1)}R)\n` +
                       `TP3: \`${signal.targets[2].price.toFixed(2)}\` (${signal.targets[2].rr.toFixed(1)}R)\n\n` +
                       `*Strength:* ${signal.strength}%\n` +
                       `*Risk:* ${signal.riskPercent.toFixed(2)}%\n` +
                       `*Best RR:* ${signal.bestRR.toFixed(2)}:1`;
        
        const url = `https://api.telegram.org/bot${TELEGRAM_CONFIG.botToken}/sendMessage`;
        
        await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                chat_id: TELEGRAM_CONFIG.chatId,
                text: message,
                parse_mode: 'Markdown'
            })
        });
        
        console.log('✅ Telegram alert sent!');
        return true;
    } catch (error) {
        console.error('❌ Telegram send failed:', error);
        return false;
    }
}

// Export functions
window.sendWhatsAppAlert = sendWhatsAppAlert;
window.sendTelegramAlert = sendTelegramAlert;
window.WHATSAPP_CONFIG = WHATSAPP_CONFIG;
window.TELEGRAM_CONFIG = TELEGRAM_CONFIG;
