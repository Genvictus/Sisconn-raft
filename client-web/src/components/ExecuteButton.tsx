
type ExecuteButtonProps = {
    onClick: () => void;
};

export default function ExecuteButton({
    onClick
}: ExecuteButtonProps) {
    return (
        <button onClick={onClick}>Execute</button>
    );
}
