import java.io.*;
import java.util.*;

// Reads T cases from stdin and prints one line per case. Each case is k, then k lists of busy intervals.
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

    static int[][] readGrid() throws IOException {
        int rows = nextInt();
        int cols = nextInt();
        int[][] g = new int[rows][cols];
        for (int r = 0; r < rows; r++) {
            for (int c = 0; c < cols; c++) {
                g[r][c] = nextInt();
            }
        }
        return g;
    }

    static String show(int[][] g) {
        StringBuilder b = new StringBuilder("[");
        for (int r = 0; r < g.length; r++) {
            if (r > 0) {
                b.append(", ");
            }
            b.append(show(g[r]));
        }
        return b.append(']').toString();
    }

    // Reads "n" then n pairs "start end".
    static int[][] readIntervals() throws IOException {
        int[][] iv = new int[nextInt()][2];
        for (int[] p : iv) {
            p[0] = nextInt();
            p[1] = nextInt();
        }
        return iv;
    }

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        br = new BufferedReader(new InputStreamReader(System.in));
        int t = nextInt();
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            int[][][] busy = new int[nextInt()][][];
            for (int i = 0; i < busy.length; i++) {
                busy[i] = readIntervals();
            }
            try {
                out.append(show(new Solution().freeTime(busy))).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
