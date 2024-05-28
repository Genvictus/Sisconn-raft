
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
        <label>
            {name}
            <input
                type="text"
                value={value}
                onChange={(e) => setValue(e.target.value)} />
        </label>
    );
}
