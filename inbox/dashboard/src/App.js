import React from 'react';
import JsonEditor from './components/JsonEditor';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>GoNetWorker</h1>
        <h3>@D7MeKz</h3>
      </header>
      <JsonEditor /> {/* JsonEditor component remains as the main content */}
    </div>
  );
}

export default App;