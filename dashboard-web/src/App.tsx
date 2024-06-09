import axios, { AxiosError } from 'axios';
import { useEffect, useState } from 'react';
import './App.css';
import DashboardTable from './components/DashboardTable';
import LogModal from './components/LogModal';
import ServerConfiguration from './components/ServerConfiguration';
import RaftNode from './types/RaftNode';

function App() {
  const [nodes, setNodes] = useState<RaftNode[]>([]);
  const [selectedLogNode, setSelectedLogNode] = useState<number>(0);
  const [isLogModalOpen, setIsLogModalOpen] = useState(false);

  const dummyNode: RaftNode = {
    Address: '',
    State: '',
    Log: ''
  }

  const openLogModal = (index: number) => {
    setSelectedLogNode(index);
    setIsLogModalOpen(true);
  }

  const closeLogModal = () => {
    setIsLogModalOpen(false);
  }

  const loadNodes = (serverAddress: string) => {
    axios.get<RaftNode[]>(`${serverAddress}/node`)
    .then((response) => {
      console.log(response.data);
      setNodes(response.data);
    }).catch((error: AxiosError) => {
      console.error("Error fetching nodes: ", error.message);
    });
  }

  const addNode = () => {
    const newNode: RaftNode = {
      Address: `http://localhost:${Math.floor(Math.random() * 10000)}`,
      State: 'Follower',
      Log: 'tes\ntest\ntest\ntes\ntest\ntest\ntes\ntest\ntest\ntes\ntest\ntest\ntes\ntest\ntest\ntes\ntest\ntest\n\ntes\ntest\ntest\n\ntes\ntest\ntest\n'
    };
    setNodes([...nodes, newNode]);
  };

  useEffect(() => {
    document.title = 'Sisconn Raft Management Dashboard';
  }, []);

  return (
    <div className="w-screen h-screen">
      <h1 className='pb-20 font-semibold'>Sisconn Raft Management Dashboard</h1>
      
      <div className="flex m-4">
        <LogModal
          modalIsOpen={isLogModalOpen}
          closeModal={closeLogModal}
          node={selectedLogNode < nodes.length ? nodes[selectedLogNode] : dummyNode}
        />

        <div className='w-2/3'>
          <DashboardTable
            nodes={nodes}
            openLogModal={openLogModal}
            setNodes={setNodes}
          />
        </div>

        <div className='w-1/3'>
          <ServerConfiguration
            loadNodes={loadNodes}
            addNode={addNode}
          />
        </div>
      </div>
    </div>
  )
}

export default App
