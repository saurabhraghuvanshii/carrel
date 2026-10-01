import java.io.*;
import java.util.*;

// Reads T cases from stdin. Each case is k calls: "new", then "put key value ttl time", "get key time" or "count time".
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
            case "new": return 0;
            case "put": return 4;
            case "get": return 2;
            case "count": return 1;
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
                TTLCache obj = null;
                for (int i = 0; i < k; i++) {
                    if (i > 0) {
                        line.append(' ');
                    }
                    switch (names[i]) {
                        case "new":
                            obj = new TTLCache();
                            line.append("null");
                            break;
                        case "put":
                            obj.put(params[i][0], params[i][1], params[i][2], params[i][3]);
                            line.append("null");
                            break;
                        case "get":
                            line.append(obj.get(params[i][0], params[i][1]));
                            break;
                        case "count":
                            line.append(obj.count(params[i][0]));
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
