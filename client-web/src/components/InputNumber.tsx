
type InputNumberProps = {
    name: string;
    value: number;
    setValue: React.Dispatch<React.SetStateAction<number>>;
};

export default function InputNumber({
    name,
    value,
    setValue
}: InputNumberProps) {
    return (
        <label>
            {name}
            <input
                type="number"
                value={value}
                onChange={(e) => setValue(parseInt(e.target.value))} />
        </label>
    );
}
