import React, { useState } from 'react';
import axios from 'axios';
import { Container, TextField, Button, Typography, Alert, Snackbar } from '@mui/material';

const App: React.FC = () => {
  const [jsonData, setJsonData] = useState<string>('');
  const [file, setFile] = useState<File | null>(null);
  const [notification, setNotification] = useState<{ message: string; success: boolean } | null>(null);

  const handleApply = async () => {
    try {
      const response = await axios.post('/api/json', JSON.parse(jsonData), {
        headers: {
          'Content-Type': 'application/json',
        },
      });
      setNotification({ message: 'JSON data sent successfully!', success: true });
      console.log('Response:', response.data);
    } catch (error) {
      setNotification({ message: 'Failed to send JSON data', success: false });
      console.error('Error applying JSON:', error);
    }
  };

  const handleUpload = async () => {
    if (!file) {
      setNotification({ message: 'No file selected', success: false });
      return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await axios.post('/api/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      setNotification({ message: 'File uploaded successfully!', success: true });
      console.log('File uploaded:', response.data);
    } catch (error) {
      setNotification({ message: 'Failed to upload file', success: false });
      console.error('Error uploading file:', error);
    }
  };

  return (
      <Container maxWidth="sm" style={{ marginTop: '50px' }}>
        <Typography variant="h4" gutterBottom>
          JSON Input and File Upload
        </Typography>
        <TextField
            label="JSON Data"
            multiline
            rows={10}
            variant="outlined"
            fullWidth
            value={jsonData}
            onChange={(e) => setJsonData(e.target.value)}
            placeholder="Enter JSON data here"
            style={{ marginBottom: '20px' }}
        />
        <Button variant="contained" color="primary" onClick={handleApply} style={{ marginBottom: '20px' }}>
          Apply
        </Button>
        <input
            accept=".json"
            style={{ display: 'none' }}
            id="upload-file"
            type="file"
            onChange={(e) => setFile(e.target.files?.[0] || null)}
        />
        <label htmlFor="upload-file">
          <Button variant="contained" component="span" color="secondary" style={{ marginBottom: '20px' }}>
            Upload File
          </Button>
        </label>
        {notification && (
            <Snackbar
                open={Boolean(notification)}
                autoHideDuration={6000}
                onClose={() => setNotification(null)}
            >
              <Alert onClose={() => setNotification(null)} severity={notification.success ? 'success' : 'error'}>
                {notification.message}
              </Alert>
            </Snackbar>
        )}
      </Container>
  );
};

export default App;
