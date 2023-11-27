using System;
using System.IO;
using System.IO.Pipes;
using System.Security.Principal;

class Program
{
    static void Main(string[] args)
    {
        using (var pipeClient = new NamedPipeClientStream(".", "temp", PipeDirection.Out)) // this pipes to: \\.\pipe\temp
        {
            pipeClient.Connect();
            using (var writer = new StreamWriter(pipeClient))
            {
                writer.AutoFlush = true;
                string currentUser = WindowsIdentity.GetCurrent().Name;
                writer.WriteLine(currentUser);
            }
        }
    }
}
