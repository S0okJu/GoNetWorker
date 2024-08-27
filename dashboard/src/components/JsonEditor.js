import React, { useState } from 'react';
import Editor from '@monaco-editor/react';
import defaultTemplate from '../data/defaultTemplate.json';

const JsonEditor = () => {
    // Initialize state with the imported JSON template
    const [jsonData, setJsonData] = useState(JSON.stringify(defaultTemplate, null, 2));
    const [error, setError] = useState('');
  
    const handleEditorChange = (value) => {
      setError('');  // Clear previous errors
      setJsonData(value);
    };
  
    const validateJson = (jsonString) => {
      try {
        JSON.parse(jsonString);
        return true;
      } catch (e) {
        return false;
      }
    };
  
    const handleSubmit = async () => {
      if (!validateJson(jsonData)) {
        setError('Invalid JSON format');
        return;
      }
  
      try {
        const response = await fetch('http://localhost:8080/data', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: jsonData,
        });
  
        if (response.ok) {
          alert('Data successfully sent!');
        } else {
          alert('Failed to send data');
        }
      } catch (e) {
        alert('Error sending data: ' + e.message);
      }
    };
  
    return (
      <div style={{ width: '90%', maxWidth: '800px', margin: '0 auto' }}>
        <Editor
          height="50vh"
          width="100%"
          defaultLanguage="json"
          value={jsonData}
          onChange={handleEditorChange}
          options={{
            formatOnType: true,
            automaticLayout: true,
            tabSize: 2,
            minimap: { enabled: false },
          }}
        />
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <button onClick={handleSubmit}>Send JSON</button>
      </div>
    );
  };
  
  export default JsonEditor;