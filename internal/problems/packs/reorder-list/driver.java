import java.io.*;
import java.util.*;

class ListNode {
    int val;
    ListNode next;
    ListNode(int val) { this.val = val; }
}

// Reads T cases from stdin and prints one line per case. Each case is a list: "n" then n values.
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

    static ListNode build(int[] vals) {
        ListNode dummy = new ListNode(0);
        ListNode tail = dummy;
        for (int v : vals) {
            tail.next = new ListNode(v);
            tail = tail.next;
        }
        return dummy.next;
    }

    // Prints a list as [a, b, c]. More than limit nodes means a loop.
    static String showList(ListNode head, int limit) {
        StringBuilder b = new StringBuilder("[");
        int count = 0;
        for (ListNode p = head; p != null; p = p.next) {
            if (++count > limit) {
                throw new IllegalStateException("the returned list has more nodes than it should. Is there a loop?");
            }
            if (count > 1) {
                b.append(", ");
            }
            b.append(p.val);
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
            try {
                ListNode head = build(vals);
                new Solution().reorder(head);
                out.append(showList(head, vals.length)).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
