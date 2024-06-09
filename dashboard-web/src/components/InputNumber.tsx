
type InputNumberProps = {
    name: string;
    value: number;
    setValue: (value: number) => void;
};

export default function InputNumber({
    name,
    value,
    setValue
}: InputNumberProps) {
    return (
        <label className="block mb-2">
            <span className="mr-2">{name}</span>
            <input
                type="number"
                value={value}
                onChange={(e) => setValue(parseInt(e.target.value))}
                className="block w-full rounded-md"
            />
        </label>
    );
}
