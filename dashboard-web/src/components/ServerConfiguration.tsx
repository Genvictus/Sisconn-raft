import React, { useMemo, useState } from "react";
import InputNumber from "./InputNumber";
import InputText from "./InputText";

type ServerConfigurationProps = {
  loadNodes: (serverAddress: string) => void;
  addNode: () => void;
};

const ServerConfiguration: React.FC<ServerConfigurationProps> = ({
  loadNodes,
  addNode,
}: ServerConfigurationProps) => {
  // Server Config
  const [serverHost, setServerHost] = useState('localhost');
  const [serverPort, setServerPort] = useState(2000);
  const [nodeHost, setNodeHost] = useState('localhost');
  const [nodePort, setNodePort] = useState(3000);
  const serverAddress = useMemo(() => `http://${serverHost}:${serverPort}`, [serverHost, serverPort]);

  return (
    <div className="pl-10">
      <h2 className="text-xl font-bold mb-4">Server Configuration</h2>
      <div className="border rounded p-4 mb-4">
        <InputText name="Server Host:" value={serverHost} setValue={setServerHost} />
        <InputNumber name="Server Port:" value={serverPort} setValue={setServerPort} />
      </div>
      

      <h2 className="text-xl font-bold mb-4">Node Configuration</h2>
      <div className="border rounded p-4 mb-4">
      <InputText name="Server Host:" value={nodeHost} setValue={setNodeHost} />
      <InputNumber name="Server Port:" value={nodePort} setValue={setNodePort} />
      <button
        className='my-5 bg-green-600 hover:bg-green-700'
        onClick={addNode}
      >
        Add New Node
      </button>
      </div>

      <button
        className='my-5 w-1/2 bg-yellow-700 hover:bg-yellow-1000'
        onClick={() => loadNodes(serverAddress)}
      >
        Fetch Nodes
      </button>
    </div>
  );
};

export default ServerConfiguration;