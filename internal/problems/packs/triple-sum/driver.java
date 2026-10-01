import java.io.*;
import java.util.*;

// Reads T cases from stdin and prints one line per case. Each case is "n", n values, then the target.
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

    // Sorts each triple, then the list, so any order is accepted.
    static String showTriples(List<List<Integer>> lists) {
        List<List<Integer>> sorted = new ArrayList<>();
        for (List<Integer> l : lists) {
            List<Integer> copy = new ArrayList<>(l);
            Collections.sort(copy);
            sorted.add(copy);
        }
        sorted.sort((a, b) -> {
            for (int i = 0; i < Math.min(a.size(), b.size()); i++) {
                int c = Integer.compare(a.get(i), b.get(i));
                if (c != 0) {
                    return c;
                }
            }
            return Integer.compare(a.size(), b.size());
        });
        StringBuilder b = new StringBuilder("[");
        for (int i = 0; i < sorted.size(); i++) {
            if (i > 0) {
                b.append(", ");
            }
            b.append('[');
            for (int j = 0; j < sorted.get(i).size(); j++) {
                if (j > 0) {
                    b.append(", ");
                }
                b.append(sorted.get(i).get(j));
            }
            b.append(']');
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
            int[] nums = readArray();
            int target = nextInt();
            try {
                out.append(showTriples(new Solution().tripleSum(nums, target))).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
