import React, { useMemo, useState } from "react";
import InputNumber from "./InputNumber";
import InputText from "./InputText";

type ServerConfigurationProps = {
  loadNodes: () => void;
  addNode: () => void;
};

const ServerConfiguration: React.FC<ServerConfigurationProps> = ({
  loadNodes,
  addNode,
}: ServerConfigurationProps) => {
  // Server Config
  const [serverHost, setServerHost] = useState('localhost');
  const [serverPort, setServerPort] = useState(2000);
  const serverAddress = useMemo(() => `http://${serverHost}:${serverPort}`, [serverHost, serverPort]);

  return (
    <div>
      <h2 className="text-xl font-bold mb-4">Server Configuration</h2>
      <div className="border rounded p-4 mb-4">
        <InputText name="Server Host:" value={serverHost} setValue={setServerHost} />
        <InputNumber name="Server Port:" value={serverPort} setValue={setServerPort} />
      </div>
      
      <button
        className='my-5 bg-green-500 hover:bg-green-700'
        onClick={loadNodes}
      >
        Fetch Nodes
      </button>
      <button
        className='my-5 bg-green-500 hover:bg-green-700'
        onClick={addNode}
      >
        Add New Node
      </button>
    </div>
  );
};

export default ServerConfiguration;