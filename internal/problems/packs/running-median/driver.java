import java.io.*;
import java.util.*;

// Reads T cases from stdin and prints one line per case. Each case is k calls: "new", "add x" or "median".
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
            int k = nextInt();
            String[] names = new String[k];
            int[] values = new int[k];
            for (int i = 0; i < k; i++) {
                names[i] = next();
                if (names[i].equals("add")) {
                    values[i] = nextInt();
                }
            }
            try {
                StringBuilder line = new StringBuilder();
                MedianFinder finder = null;
                for (int i = 0; i < k; i++) {
                    if (i > 0) {
                        line.append(' ');
                    }
                    switch (names[i]) {
                        case "new":
                            finder = new MedianFinder();
                            line.append("null");
                            break;
                        case "add":
                            finder.add(values[i]);
                            line.append("null");
                            break;
                        default:
                            line.append(String.format(Locale.ROOT, "%.1f", finder.median()));
                    }
                }
                out.append(line).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
