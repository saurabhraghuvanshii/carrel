import java.io.*;
import java.util.*;

class ListNode {
    int val;
    ListNode next;
    ListNode(int val) { this.val = val; }
}

// Reads T cases from stdin and prints one line per case. Each case is a list, then the position the last node points back to, or -1.
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

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        br = new BufferedReader(new InputStreamReader(System.in));
        int t = nextInt();
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            int[] vals = readArray();
            int pos = nextInt();
            ListNode[] nodes = new ListNode[vals.length];
            for (int i = 0; i < vals.length; i++) {
                nodes[i] = new ListNode(vals[i]);
                if (i > 0) {
                    nodes[i - 1].next = nodes[i];
                }
            }
            if (pos >= 0) {
                nodes[vals.length - 1].next = nodes[pos];
            }
            ListNode head = vals.length == 0 ? null : nodes[0];
            try {
                ListNode start = new Solution().loopStart(head);
                int at = -1;
                for (int i = 0; i < nodes.length; i++) {
                    if (nodes[i] == start) {
                        at = i;
                    }
                }
                if (start != null && at < 0) {
                    out.append("ERROR returned a node that is not in the list\n");
                } else {
                    out.append(at).append('\n');
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
