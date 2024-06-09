import React, { useState } from 'react';
import Node from '../types/Node';

const DashboardTable: React.FC = () => {
  const [nodes, setNodes] = useState<Node[]>([]);

  const addNode = () => {
    const newNode: Node = {
      address: `http://localhost:${Math.floor(Math.random() * 10000)}`,
      state: 'Follower',
      log: 'tes\ntest\ntest\n'
    };
    setNodes([...nodes, newNode]);
  };

  const deleteNode = (addr: string) => {
    setNodes(nodes.filter(node => node.address !== addr));
  };

  const handleLogChange = (addr: string, value: string) => {
    setNodes(nodes.map(node => node.address === addr ? { ...node, log: value } : node));
  };

  return (
    <div>
      <h1>Sisconn Raft Management Dashboard</h1>
      <button className='my-5 bg-green-500 hover:bg-green-700' onClick={addNode}>Add New Node</button>
      <table className='w-full'>
        <thead>
          <tr>
            <th>No</th>
            <th>State</th>
            <th>Address</th>
            <th>Log</th>
            <th>Delete</th>
          </tr>
        </thead>
        <tbody>
          {nodes.map((node, index) => (
            <tr key={node.address}>
              <td>{index + 1}</td>
              <td>{node.state}</td>
              <td>{node.address}</td>
              <td className=''>
                <button className=' text-sm bg-blue-500 hover:bg-blue-700'>View Log</button>
              </td>
              <td className=''>
                <button className=' text-sm bg-red-500 hover:bg-red-700' onClick={() => deleteNode(node.address)}>Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default DashboardTable;
