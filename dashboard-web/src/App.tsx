import axios, { AxiosError } from 'axios';
import { useEffect, useMemo, useState } from 'react';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './App.css';
import DashboardTable from './components/DashboardTable';
import LogModal from './components/LogModal';
import ServerConfiguration from './components/ServerConfiguration';
import RaftNode from './types/RaftNode';
import { splitAddress } from './utils/util';

function App() {
  const [nodes, setNodes] = useState<RaftNode[]>([]);
  const [selectedLogNode, setSelectedLogNode] = useState<number>(0);
  const [isLogModalOpen, setIsLogModalOpen] = useState(false);

  // Server Config
  const [serverHost, setServerHost] = useState('localhost');
  const [serverPort, setServerPort] = useState(2000);
  const serverAddress = useMemo(() => `http://${serverHost}:${serverPort}`, [serverHost, serverPort]);

  const dummyNode: RaftNode = {
    Address: '',
    State: '',
    Log: ''
  }

  const notifyToast = (message: string, success: boolean) => {
    if (success) {
      toast.success(message);
    } else {
      toast.error(message);
    }
  }

  const openLogModal = (index: number) => {
    const { host, port } = splitAddress(nodes[index].Address);
    axios.get<string>(`${serverAddress}/log`, { params: { host, port } })
    .then((response) => {
        setNodes(nodes.map((node, idx) => idx === index ? { ...node, Log: response.data } : node));
      }).catch((error: AxiosError) => {
        console.error("Error fetching logs: ", error.message);
        notifyToast('Error fetching logs', false);
      });
  
    setSelectedLogNode(index);
    setIsLogModalOpen(true);
  }

  const closeLogModal = () => {
    setIsLogModalOpen(false);
  }

  const loadNodes = () => {
    axios.get<RaftNode[]>(`${serverAddress}/node`)
      .then((response) => {
        setNodes(response.data);
        notifyToast('Nodes loaded successfully', true);
      }).catch((error: AxiosError) => {
        console.error("Error fetching nodes: ", error.message);
        notifyToast('Error fetching nodes', false);
      });
  }

  const addNode = (host: string, port: number) => {
    axios.get(`${serverAddress}/node/add`, { params: { host, port } })
    .then(() => {
      const newNode: RaftNode = {
        Address: `${host}:${port}`,
        State: 'FOLLOWER',
        Log: ''
      };
      setNodes([...nodes, newNode]);
      notifyToast('Nodes added successfully', true);
    }).catch((error: AxiosError) => {
      console.error("Error adding nodes: ", error.message);
      notifyToast('Error adding nodes', false);
    });
  };
  
  const deleteNode = (index: number) => {
    const { host, port } = splitAddress(nodes[index].Address);
    axios.get(`${serverAddress}/node/delete`, { params: { host, port } })
    .then(() => {
      const newNode: RaftNode = {
        Address: `${host}:${port}`,
        State: 'FOLLOWER',
        Log: ''
      };
      setNodes(nodes.filter((_, idx) => idx !== index));
      notifyToast('Remove Nodes successfully', true);
    }).catch((error: AxiosError) => {
      console.error("Error removing nodes: ", error.message);
      notifyToast('Error removing nodes', false);
    });
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
            deleteNode={deleteNode}
          />
        </div>

        <div className='w-1/3'>
          <ServerConfiguration
            serverHost={serverHost}
            serverPort={serverPort}
            setServerHost={setServerHost}
            setServerPort={setServerPort}
            loadNodes={loadNodes}
            addNode={addNode}
          />
        </div>
      </div>
      
      <ToastContainer
        position="bottom-right"
        autoClose={5000}
        hideProgressBar={false}
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
      />
    </div>
  )
}

export default App
