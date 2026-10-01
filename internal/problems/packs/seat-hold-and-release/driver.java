import java.io.*;
import java.util.*;

// Reads T cases from stdin. Each case is k calls: "new seats ttl", then "hold user count time", "buy user time", "release user time" or "freeSeats time".
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
            case "new": return 2;
            case "hold": return 3;
            case "buy": return 2;
            case "release": return 2;
            case "freeSeats": return 1;
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
                SeatHolder obj = null;
                for (int i = 0; i < k; i++) {
                    if (i > 0) {
                        line.append(' ');
                    }
                    switch (names[i]) {
                        case "new":
                            obj = new SeatHolder(params[i][0], params[i][1]);
                            line.append("null");
                            break;
                        case "hold":
                            line.append(obj.hold(params[i][0], params[i][1], params[i][2]));
                            break;
                        case "buy":
                            line.append(obj.buy(params[i][0], params[i][1]));
                            break;
                        case "release":
                            line.append(obj.release(params[i][0], params[i][1]));
                            break;
                        case "freeSeats":
                            line.append(obj.freeSeats(params[i][0]));
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
