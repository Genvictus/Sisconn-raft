import React from 'react';
import RaftNode from '../types/RaftNode';

type DashboardTableProps = {
  nodes: RaftNode[]
  openLogModal: (index: number) => void
  setNodes: (nodes: RaftNode[]) => void
};

const DashboardTable: React.FC<DashboardTableProps> = ({
  nodes,
  openLogModal,
  setNodes,
}: DashboardTableProps) => {

  const deleteNode = (index: number) => {
    setNodes(nodes.filter((_, idx) => idx !== index));
  };

  const handleLogChange = (index: number, value: string) => {
    setNodes(nodes.map((node, idx) => idx === index ? { ...node, log: value } : node));
  };

  return (
    <>
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
        <tbody >
          {nodes.map((node, index) => (
            <tr key={index}>
              <td>{index + 1}</td>
              <td>{node.State}</td>
              <td>{node.Address}</td>
              <td>
                <button
                  className=' m-2 text-sm bg-blue-600 hover:bg-blue-700'
                  onClick={() => openLogModal(index)}
                >
                  View Log
                </button>
              </td>
              <td>
                <button
                  className=' text-sm bg-red-600 hover:bg-red-700'
                  onClick={() => deleteNode(index)}
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
};

export default DashboardTable;
