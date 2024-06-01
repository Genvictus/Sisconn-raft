import { useEffect, useRef } from "react";

type LogProps = {
    logs: string;
};

export default function Log({
    logs
}: LogProps) {
    const textAreaRef = useRef<HTMLTextAreaElement | null>(null);

    useEffect(() => {
        const textArea = textAreaRef.current;
        if (textArea) {
            textArea.scrollTop = textArea.scrollHeight;
        }
    }, [logs]);

    return (
        <div className="h-full w-full flex justify-center items-center p-4">
            <div className="h-full w-full border border-gray-300 rounded-md overflow-hidden">
                <textarea
                    className="w-full h-full p-4 resize-none overflow-auto outline-none border-none"
                    value={logs}
                    ref={textAreaRef}
                    readOnly >
                </textarea>
            </div>
        </div>
    );
}