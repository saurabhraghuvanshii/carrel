import java.io.*;
import java.util.*;

// Reads T cases from stdin and prints one line per case. Each case is "n q", then q operations: "get i", "set i", "clear i" or "update i b".
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints goes to stderr so it cannot break the results.
public class Main {
    static BufferedReader br;
    static StringTokenizer st;

    static String next() throws IOException {
        while (st == null || !st.hasMoreTokens()) {
            st = new StringTokenizer(br.readLine());
        }
        return st.nextToken();
    }

    static int nextInt() throws IOException {
        return Integer.parseInt(next());
    }

    static int[] readArray() throws IOException {
        int[] a = new int[nextInt()];
        for (int i = 0; i < a.length; i++) {
            a[i] = nextInt();
        }
        return a;
    }

    static String show(int[] a) {
        StringBuilder b = new StringBuilder("[");
        for (int i = 0; i < a.length; i++) {
            if (i > 0) {
                b.append(", ");
            }
            b.append(a[i]);
        }
        return b.append(']').toString();
    }

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        br = new BufferedReader(new InputStreamReader(System.in));
        int t = nextInt();
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            int n = nextInt();
            int q = nextInt();
            String[] ops = new String[q];
            int[] at = new int[q];
            int[] bit = new int[q];
            for (int k = 0; k < q; k++) {
                ops[k] = next();
                at[k] = nextInt();
                if (ops[k].equals("update")) {
                    bit[k] = nextInt();
                }
            }
            try {
                Solution sol = new Solution();
                int[] results = new int[q];
                for (int k = 0; k < q; k++) {
                    if (ops[k].equals("get")) {
                        results[k] = sol.getBit(n, at[k]);
                        continue;
                    }
                    if (ops[k].equals("set")) {
                        n = sol.setBit(n, at[k]);
                    } else if (ops[k].equals("clear")) {
                        n = sol.clearBit(n, at[k]);
                    } else {
                        n = sol.updateBit(n, at[k], bit[k]);
                    }
                    results[k] = n;
                }
                out.append(show(results)).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
