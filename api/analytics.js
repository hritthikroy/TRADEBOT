// Vercel Serverless Function for Analytics
const { createClient } = require('@supabase/supabase-js');

const supabase = createClient(
  'https://xlxugbqxfrrwutxecwug.supabase.co',
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InhseHVnYnF4ZnJyd3V0eGVjd3VnIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ0MTM0MDQsImV4cCI6MjA3OTk4OTQwNH0.LAQKcsWOOTOeAcKh9zRF9l740aDp0Ki1cws8oG3yQZU'
);

module.exports = async function handler(req, res) {
  res.setHeader('Access-Control-Allow-Origin', '*');
  
  try {
    const { data, error } = await supabase
      .from('signal_analytics')
      .select('*');

    if (error) throw error;
    return res.status(200).json(data);
  } catch (error) {
    return res.status(500).json({ error: error.message });
  }
}
