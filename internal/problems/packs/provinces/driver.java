import java.io.*;
import java.util.*;

// Reads T cases from stdin and prints one line per case. Each case is "n m", then m edges "u v".
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

    // Reads m edges of width numbers each: "u v" or "u v w".
    static int[][] readEdges(int m, int width) throws IOException {
        int[][] edges = new int[m][width];
        for (int[] e : edges) {
            for (int j = 0; j < width; j++) {
                e[j] = nextInt();
            }
        }
        return edges;
    }

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        br = new BufferedReader(new InputStreamReader(System.in));
        int t = nextInt();
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            int n = nextInt();
            int m = nextInt();
            int[][] edges = readEdges(m, 2);
            try {
                out.append(new Solution().countGroups(n, edges)).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
