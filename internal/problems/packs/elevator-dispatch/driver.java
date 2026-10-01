import java.io.*;
import java.util.*;

// Reads T cases from stdin. Each case is k calls: "new floors", then "call floor", "next" or "travelled".
// Prints one line per case with one result per call, null for calls that
// return nothing. Flushes after every case so a crash cannot hide finished
// results. Anything the learner prints goes to stderr.
public class Main {
    static BufferedReader br;
    static StringTokenizer st;

    static String next() throws IOException {
        while (st == null || !st.hasMoreTokens()) {
            st = new StringTokenizer(br.readLine());
        }
        return st.nextToken();
    }

    static int argCount(String name) {
        switch (name) {
            case "new": return 1;
            case "call": return 1;
            case "next": return 0;
            case "travelled": return 0;
            default: throw new IllegalArgumentException("unknown call " + name);
        }
    }

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        br = new BufferedReader(new InputStreamReader(System.in));
        int t = Integer.parseInt(next());
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            int k = Integer.parseInt(next());
            String[] names = new String[k];
            int[][] params = new int[k][];
            for (int i = 0; i < k; i++) {
                names[i] = next();
                params[i] = new int[argCount(names[i])];
                for (int j = 0; j < params[i].length; j++) {
                    params[i][j] = Integer.parseInt(next());
                }
            }
            StringBuilder line = new StringBuilder();
            try {
                Elevator obj = null;
                for (int i = 0; i < k; i++) {
                    if (i > 0) {
                        line.append(' ');
                    }
                    switch (names[i]) {
                        case "new":
                            obj = new Elevator(params[i][0]);
                            line.append("null");
                            break;
                        case "call":
                            obj.call(params[i][0]);
                            line.append("null");
                            break;
                        case "next":
                            line.append(obj.next());
                            break;
                        case "travelled":
                            line.append(obj.travelled());
                            break;
                    }
                }
            } catch (Throwable e) {
                line = new StringBuilder("ERROR ").append(e);
            }
            out.append(line).append('\n');
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
