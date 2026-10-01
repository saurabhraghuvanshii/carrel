import java.io.*;
import java.util.*;

// Flushes after every case so a crash cannot hide finished results.
// Reads T cases from stdin, calls the learner's code, prints one line per case.
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

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        br = new BufferedReader(new InputStreamReader(System.in));
        int t = Integer.parseInt(next());
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            int n = Integer.parseInt(next());
            int[] nums = new int[n];
            for (int i = 0; i < n; i++) {
                nums[i] = Integer.parseInt(next());
            }
            int target = Integer.parseInt(next());
            try {
                int[] r = new Solution().pairWithTarget(nums, target);
                if (r == null || r.length != 2) {
                    out.append("ERROR bad return value\n");
                } else {
                    out.append(Math.min(r[0], r[1])).append(' ').append(Math.max(r[0], r[1])).append('\n');
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
