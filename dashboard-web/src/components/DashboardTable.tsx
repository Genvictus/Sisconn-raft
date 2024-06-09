import React, { useState } from 'react';
import Node from '../types/Node';
import LogModal from './LogModal';

const DashboardTable: React.FC = () => {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [isLogModalOpen, setIsLogModalOpen] = useState(false);
  const [selectedLogNode, setSelectedLogNode] = useState<number>(0);

  const dummyNode: Node = {
    address: '',
    state: '',
    log: ''
  }

  const openLogModal = (index: number) => {
    setSelectedLogNode(index);
    setIsLogModalOpen(true);
  }

  const closeLogModal = () => {
    setIsLogModalOpen(false);
  }

  const addNode = () => {
    const newNode: Node = {
      address: `http://localhost:${Math.floor(Math.random() * 10000)}`,
      state: 'Follower',
      log: 'tes\ntest\ntest\ntes\ntest\ntest\ntes\ntest\ntest\ntes\ntest\ntest\ntes\ntest\ntest\ntes\ntest\ntest\n\ntes\ntest\ntest\n\ntes\ntest\ntest\n'
    };
    setNodes([...nodes, newNode]);
  };

  const deleteNode = (index: number) => {
    setNodes(nodes.filter((_, idx) => idx !== index));
  };

  const handleLogChange = (index: number, value: string) => {
    setNodes(nodes.map((node, idx) => idx === index ? { ...node, log: value } : node));
  };

  return (
    <div>
      <h1>Sisconn Raft Management Dashboard</h1>
      <button
        className='my-5 bg-green-500 hover:bg-green-700'
        onClick={addNode}
      >
        Add New Node
      </button>

      <LogModal
        modalIsOpen={isLogModalOpen}
        closeModal={closeLogModal}
        node={selectedLogNode < nodes.length ? nodes[selectedLogNode] : dummyNode}
      />

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
              <td>
                <button
                  className=' text-sm bg-blue-500 hover:bg-blue-700'
                  onClick={() => openLogModal(index)}
                >
                  View Log
                </button>
              </td>
              <td>
                <button
                  className=' text-sm bg-red-500 hover:bg-red-700'
                  onClick={() => deleteNode(index)}
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default DashboardTable;
