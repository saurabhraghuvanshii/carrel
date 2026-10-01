import java.io.*;
import java.util.*;

class ListNode {
    int val;
    ListNode next;
    ListNode(int val) { this.val = val; }
}

// Flushes after every case so a crash cannot hide finished results.
// Reads T cases from stdin, each "n" then n values. Prints the returned list
// as [a, b, c] on one line. Anything the learner prints goes to stderr.
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
            ListNode dummy = new ListNode(0);
            ListNode tail = dummy;
            for (int i = 0; i < n; i++) {
                tail.next = new ListNode(nextInt());
                tail = tail.next;
            }
            try {
                ListNode r = new Solution().reverseList(dummy.next);
                StringBuilder line = new StringBuilder("[");
                int count = 0;
                for (ListNode p = r; p != null; p = p.next) {
                    if (++count > n) {
                        throw new IllegalStateException("the returned list has more nodes than the input. Is there a loop?");
                    }
                    if (count > 1) {
                        line.append(", ");
                    }
                    line.append(p.val);
                }
                out.append(line).append("]\n");
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
