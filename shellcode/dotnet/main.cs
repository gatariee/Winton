using System;
using System.IO;
using System.IO.Pipes;

class Program
{
    static void Main(string[] args)
    {
        using (var pipeClient = new NamedPipeClientStream(".", "temp", PipeDirection.Out))
        {
            pipeClient.Connect();
            using (var writer = new StreamWriter(pipeClient))
            {
                writer.AutoFlush = true;
                writer.WriteLine("weenton");
            }
        }
    }
}
