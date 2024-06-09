
export const splitAddress = (address: string): { host: string, port: number } => {
  const [host, port] = address.split(':');
  return { host, port: parseInt(port) };
}
