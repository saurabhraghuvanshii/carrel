import java.io.*;
import java.util.*;

// Flushes after every case so a crash cannot hide finished results.
// Reads T cases from stdin, calls the learner's code, prints one line per case.
// Anything the learner prints goes to stderr so it cannot break the results.
public class Main {
    static BufferedReader br;
    static StringTokenizer st;

    static int nextInt() throws IOException {
        while (st == null || !st.hasMoreTokens()) {
            st = new StringTokenizer(br.readLine());
        }
        return Integer.parseInt(st.nextToken());
    }

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        br = new BufferedReader(new InputStreamReader(System.in));
        int t = nextInt();
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            int n = nextInt();
            int[][] drivers = new int[n][4];
            for (int i = 0; i < n; i++) {
                for (int j = 0; j < 4; j++) {
                    drivers[i][j] = nextInt();
                }
            }
            int rx = nextInt();
            int ry = nextInt();
            try {
                out.append(new Solution().nearestDriver(drivers, rx, ry)).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
