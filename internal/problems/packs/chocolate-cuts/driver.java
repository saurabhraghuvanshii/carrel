import java.io.*;
import java.util.*;

// Reads T cases from stdin and prints one line per case. Each case is "rows cols", then rows − 1 horizontal costs, then cols − 1 vertical costs.
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

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        br = new BufferedReader(new InputStreamReader(System.in));
        int t = nextInt();
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            int rows = nextInt();
            int cols = nextInt();
            int[] horizontal = new int[rows - 1];
            for (int i = 0; i < horizontal.length; i++) {
                horizontal[i] = nextInt();
            }
            int[] vertical = new int[cols - 1];
            for (int i = 0; i < vertical.length; i++) {
                vertical[i] = nextInt();
            }
            try {
                out.append(new Solution().cutCost(rows, cols, horizontal, vertical)).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
