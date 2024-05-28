import { useState } from "react";
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
    const [serverPort, setServerPort] = useState(6969);

    // COmmands
    const [getKey, setGetKey] = useState('');
    const [setKey, setSetKey] = useState('');
    const [setValue, setSetValue] = useState('');
    const [strlnKey, setStrlnKey] = useState('');
    const [delKey, setDelKey] = useState('');
    const [appendKey, setAppendKey] = useState('');
    const [appendValue, setAppendValue] = useState('');
    
    const handlePing = () => {
        // TODO
        appendLog("ping");
        setGetKey("qweqw");
    };

    const handleGet = () => {
        // TODO
    };

    const handleSet = () => {
        // TODO
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
        <div className="h-full p-4">
            <h2>Server Configuration</h2>
            
            <br />

            <InputText name="Server Host:" value={serverHost} setValue={setServerHost} />
            <br />
            <InputNumber name="Server Port:" value={serverPort} setValue={setServerPort} />

            <br />
            <br />

            <h2>Commands</h2>
            <h3></h3>
            <h3>ping</h3>
            <ExecuteButton onClick={handlePing} />

            <h3>get</h3>
            <InputText name="Key" value={getKey} setValue={setGetKey} />
            <br />
            <ExecuteButton onClick={handleGet} />
            
            <h3>set</h3>
            <InputText name="Key" value={setKey} setValue={setSetKey} />
            <br />
            <InputText name="Value" value={setValue} setValue={setSetValue} />
            <br />
            <ExecuteButton onClick={handleSet} />

            <h3>strln</h3>
            <InputText name="Key" value={strlnKey} setValue={setStrlnKey} />
            <br />
            <ExecuteButton onClick={handleStrln} />
            
            <h3>del</h3>
            <InputText name="Key" value={delKey} setValue={setDelKey} />
            <br />
            <ExecuteButton onClick={handleDel} />
            
            <h3>append</h3>
            <InputText name="Key" value={appendKey} setValue={setAppendKey} />
            <br />
            <InputText name="Value" value={appendValue} setValue={setAppendValue} />
            <br />
            <ExecuteButton onClick={handleAppend} />

            <h3>request log</h3>
            <ExecuteButton onClick={handleRequestLog} />

        </div>
    );
}