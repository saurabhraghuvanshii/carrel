import java.io.*;
import java.util.*;

// Reads T cases from stdin and prints one line per case. Each case is q, then q operations: "insert w", "search w" or "prefix p".
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
            int q = nextInt();
            String[] ops = new String[q];
            String[] arg = new String[q];
            for (int k = 0; k < q; k++) {
                ops[k] = next();
                arg[k] = next();
            }
            try {
                Trie trie = new Trie();
                List<Boolean> results = new ArrayList<>();
                for (int k = 0; k < q; k++) {
                    if (ops[k].equals("insert")) {
                        trie.insert(arg[k]);
                    } else if (ops[k].equals("search")) {
                        results.add(trie.search(arg[k]));
                    } else {
                        results.add(trie.startsWith(arg[k]));
                    }
                }
                out.append(results).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
