import axios, { AxiosError } from "axios";
import { useEffect, useMemo, useState } from "react";
import ExecuteButton from "./ExecuteButton";
import InputNumber from "./InputNumber";
import InputText from "./InputText";

type APIManagerProps = {
    appendLog: (log: string) => void;
};

export default function APIManager({
    appendLog
}: APIManagerProps) {
    // Server Config
    const [serverHost, setServerHost] = useState('localhost');
    const [serverPort, setServerPort] = useState(2000);
    const serverAddress = useMemo(() => `http://${serverHost}:${serverPort}`, [serverHost, serverPort]);

    // COmmands
    const [getKey, setGetKey] = useState('');
    const [setKey, setSetKey] = useState('');
    const [setValue, setSetValue] = useState('');
    const [strlnKey, setStrlnKey] = useState('');
    const [delKey, setDelKey] = useState('');
    const [appendKey, setAppendKey] = useState('');
    const [appendValue, setAppendValue] = useState('');

    useEffect(() => {
        handlePing()
    }, []);
    
    const handlePing = async () => {
        // TODO handle all async appendLog issue
        appendLog(`Pinging ${serverAddress}`);
        axios.get(`${serverAddress}/ping`).then((response) => {
            appendLog(response.data);
        }).catch((error: AxiosError) => {
            appendLog(error.message);
        });
    };

    const handleGet = () => {
        // TODO
    };

    const handleSet = () => {
        appendLog(`Set ${setKey} ${setValue}`);
        axios.get(`${serverAddress}/set`, {
            params: {
                key: setKey,
                value: setValue
            }
        }).then((response) => {
            appendLog(response.data);
        }).catch((error: AxiosError) => {
            appendLog(error.message);
        });

    };

    const handleStrln = () => {
        // TODO
    };

    const handleDel = () => {
        // TODO
    };

    const handleAppend = () => {
        // TODO
    };

    const handleRequestLog = () => {
        // TODO
    };

    return (
        <div className="container mx-auto p-4">
        <h2 className="text-xl font-bold mb-4">Server Configuration</h2>
    
        <div className="border rounded p-4 mb-4">
            <InputText name="Server Host:" value={serverHost} setValue={setServerHost} />
            <InputNumber name="Server Port:" value={serverPort} setValue={setServerPort} />
        </div>
    
        <h2 className="text-xl font-bold mb-2">Commands</h2>
    
        <div className="grid grid-cols-3 gap-4">
            <div className="border rounded p-4">
                <h3 className="font-medium mb-2">ping</h3>
                <ExecuteButton onClick={handlePing} />
            </div>

            <div className="border rounded p-4">
                <h3 className="font-medium mb-2">get</h3>
                <InputText name="Key" value={getKey} setValue={setGetKey} />
                <ExecuteButton onClick={handleGet} />
            </div>

            <div className="border rounded p-4">
                <h3 className="font-medium mb-2">set</h3>
                <InputText name="Key" value={setKey} setValue={setSetKey} />
                <InputText name="Value" value={setValue} setValue={setSetValue} />
                <ExecuteButton onClick={handleSet} />
            </div>

            <div className="border rounded p-4">
                <h3 className="font-medium mb-2">strln</h3>
                <InputText name="Key" value={strlnKey} setValue={setStrlnKey} />
                <ExecuteButton onClick={handleStrln} />
            </div>

            <div className="border rounded p-4">
                <h3 className="font-medium mb-2">del</h3>
                <InputText name="Key" value={delKey} setValue={setDelKey} />
                <ExecuteButton onClick={handleDel} />
            </div>

            <div className="border rounded p-4">
                <h3 className="font-medium mb-2">append</h3>
                <InputText name="Key" value={appendKey} setValue={setAppendKey} />
                <InputText name="Value" value={appendValue} setValue={setAppendValue} />
                <ExecuteButton onClick={handleAppend} />
            </div>
        </div>
    
        <div className="border rounded p-4 mt-4">
            <h3 className="font-medium mb-2">request log</h3>
            <ExecuteButton onClick={handleRequestLog} />
        </div>
    </div>
    );
}