
type InputTextProps = {
    name: string;
    value: string;
    setValue: React.Dispatch<React.SetStateAction<string>>;
};

export default function InputText({
    name,
    value,
    setValue
}: InputTextProps) {
    return (
        <label className="block mb-2">
            <span className="mr-2">{name}</span>
            <input
                type="text"
                value={value}
                onChange={(e) => setValue(e.target.value)}
                className="block w-full rounded-md"
            />
        </label>
    );
}
