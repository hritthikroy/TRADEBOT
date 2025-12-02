// Vercel Serverless Function for Signals
const { createClient } = require('@supabase/supabase-js');

const supabase = createClient(
  'https://xlxugbqxfrrwutxecwug.supabase.co',
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InhseHVnYnF4ZnJyd3V0eGVjd3VnIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ0MTM0MDQsImV4cCI6MjA3OTk4OTQwNH0.LAQKcsWOOTOeAcKh9zRF9l740aDp0Ki1cws8oG3yQZU'
);

module.exports = async function handler(req, res) {
  // Enable CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  if (req.method === 'OPTIONS') {
    return res.status(200).end();
  }

  try {
    if (req.method === 'GET') {
      // Get all signals
      const { data, error } = await supabase
        .from('trading_signals')
        .select('*')
        .order('created_at', { ascending: false });

      if (error) throw error;
      return res.status(200).json(data);
    }

    if (req.method === 'POST') {
      // Create signal
      const { data, error } = await supabase
        .from('trading_signals')
        .insert([req.body])
        .select();

      if (error) throw error;
      return res.status(201).json(data[0]);
    }

    return res.status(405).json({ error: 'Method not allowed' });
  } catch (error) {
    return res.status(500).json({ error: error.message });
  }
}
