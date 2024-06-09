import React from "react";
import ReactModal from "react-modal";
import Node from "../types/Node";

type LogModalProps = {
    modalIsOpen: boolean;
    closeModal: () => void;
    node: Node;
};

const LogModal: React.FC<LogModalProps> = ({
    modalIsOpen,
    closeModal,
    node
}: LogModalProps) => {
    return (
        <ReactModal
            isOpen={modalIsOpen}
            onRequestClose={closeModal}
            contentLabel="Log Modal"
            className="absolute overflow-auto left-20 right-20 top-10 bottom-10"
            overlayClassName="fixed inset-0 bg-gray-500 bg-opacity-50"
        >
            <div className=" bg-black p-10 rounded-lg h-full">
                <h1 className="text-2xl font-semibold mb-4">Node {node.address} Log</h1>
                <textarea
                    className='w-full h-3/4 border border-gray-300 rounded-lg px-3 py-2 mb-4 resize-none overflow-auto outline-none border-none text-xl'
                    // className="w-full h-1/2 p-4 resize-none overflow-auto outline-none border-none"
                    value={node.log}
                    readOnly >
                </textarea>
                <button
                    className='bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded text-xl'
                    onClick={closeModal}
                >
                    Close
                </button>
            </div>
        </ReactModal>
    );
};

export default LogModal;