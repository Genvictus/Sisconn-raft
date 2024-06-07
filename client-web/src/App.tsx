import { useState } from 'react';
import APIManager from './components/APIManager';
import Log from './components/Log';

function App() {
  const [logs, setLogs] = useState('');

  const appendLog = (log: string) => {
    setLogs(prevLogs => prevLogs + '> ' + log + '\n');
  };

  return (
    <div className="w-screen h-screen">
      <h1 className="text-2xl font-bold text-center pt-4 pb-8">
        Sisconn Raft Client
      </h1>

      <div className="flex">
        <div className='w-1/2'>
          <APIManager
            appendLog={appendLog}
          />
        </div>
        <div className='w-1/2'>
          <Log
            logs={logs}
          />
        </div>
      </div>
    </div>
  )
}

export default App
