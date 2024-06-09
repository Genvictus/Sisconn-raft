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
    </>
  );
};

export default DashboardTable;
