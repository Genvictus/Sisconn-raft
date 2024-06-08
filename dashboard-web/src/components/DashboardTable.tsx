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
        <button onClick={addNode}>Add New Node</button>
        <table>
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
                <td>
                  <textarea
                    disabled
                    value={node.log}
                    onChange={(e) => handleLogChange(node.address, e.target.value)}
                  />
                </td>
                <td>
                  <button onClick={() => deleteNode(node.address)}>Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  };
  
  export default DashboardTable;
