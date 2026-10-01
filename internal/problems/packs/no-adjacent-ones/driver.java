import java.io.*;
import java.util.*;

// Reads T cases from stdin and prints one line per case. Each case is n.
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

    static String quoted(List<String> words) {
        StringBuilder b = new StringBuilder("[");
        for (int i = 0; i < words.size(); i++) {
            if (i > 0) {
                b.append(", ");
            }
            b.append('"').append(words.get(i)).append('"');
        }
        return b.append(']').toString();
    }

    // Sorts the words, so any order is accepted, then prints ["a", "b"].
    static String showWords(List<String> words) {
        List<String> s = new ArrayList<>(words);
        Collections.sort(s);
        return quoted(s);
    }

    // Sorts the lists (not the words inside them), then prints [["a", "b"], ["ab"]].
    static String showSplits(List<List<String>> lists) {
        List<List<String>> s = new ArrayList<>(lists);
        s.sort((a, b) -> {
            for (int i = 0; i < Math.min(a.size(), b.size()); i++) {
                int c = a.get(i).compareTo(b.get(i));
                if (c != 0) {
                    return c;
                }
            }
            return Integer.compare(a.size(), b.size());
        });
        StringBuilder b = new StringBuilder("[");
        for (int i = 0; i < s.size(); i++) {
            if (i > 0) {
                b.append(", ");
            }
            b.append(quoted(s.get(i)));
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
            try {
                out.append(showWords(new Solution().noAdjacentOnes(n))).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
