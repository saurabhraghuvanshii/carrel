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
            int[][] before = readEdges(m, 2);
            try {
                int[] got = new Solution().order(n, before);
                if (got == null || got.length == 0) {
                    out.append("[]\n");
                } else {
                    int[] pos = new int[n];
                    Arrays.fill(pos, -1);
                    String problem = got.length != n ? "it has " + got.length + " courses, not " + n : null;
                    for (int i = 0; problem == null && i < n; i++) {
                        if (got[i] < 0 || got[i] >= n || pos[got[i]] >= 0) {
                            problem = "course " + got[i] + " is out of range or appears twice";
                        } else {
                            pos[got[i]] = i;
                        }
                    }
                    for (int[] e : before) {
                        if (problem == null && pos[e[0]] > pos[e[1]]) {
                            problem = "course " + e[1] + " comes before course " + e[0] + ", which it needs";
                        }
                    }
                    out.append(problem == null ? "valid order" : "not a valid order: " + problem).append('\n');
                }
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
