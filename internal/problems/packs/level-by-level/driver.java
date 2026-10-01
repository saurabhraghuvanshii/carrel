import java.io.*;
import java.util.*;

class TreeNode {
    int val;
    TreeNode left, right;
    TreeNode(int val) { this.val = val; }
}

// Reads T cases from stdin and prints one line per case. Each case is one tree in level order, "null" for a missing child.
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

    // Reads a tree in level order, "null" for a missing child.
    static TreeNode parseTree(String line) {
        String[] tok = line.trim().split("\\s+");
        if (tok[0].isEmpty() || tok[0].equals("null")) {
            return null;
        }
        TreeNode root = new TreeNode(Integer.parseInt(tok[0]));
        ArrayDeque<TreeNode> queue = new ArrayDeque<>();
        queue.add(root);
        for (int i = 1; i < tok.length; ) {
            TreeNode node = queue.poll();
            if (!tok[i].equals("null")) {
                node.left = new TreeNode(Integer.parseInt(tok[i]));
                queue.add(node.left);
            }
            i++;
            if (i < tok.length && !tok[i].equals("null")) {
                node.right = new TreeNode(Integer.parseInt(tok[i]));
                queue.add(node.right);
            }
            i++;
        }
        return root;
    }

    static int countNodes(String line) {
        int n = 0;
        for (String tok : line.trim().split("\\s+")) {
            if (!tok.isEmpty() && !tok.equals("null")) {
                n++;
            }
        }
        return n;
    }

    // Prints a tree in level order as [a, b, null, c]. More than limit nodes
    // means the tree has a loop.
    static String showTree(TreeNode root, int limit) {
        if (root == null) {
            return "[]";
        }
        List<String> parts = new ArrayList<>();
        LinkedList<TreeNode> queue = new LinkedList<>();
        queue.add(root);
        int count = 0;
        while (!queue.isEmpty()) {
            TreeNode n = queue.poll();
            if (n == null) {
                parts.add("null");
                continue;
            }
            if (++count > limit) {
                throw new IllegalStateException("the returned tree has more nodes than it should. Is there a loop?");
            }
            parts.add(String.valueOf(n.val));
            queue.add(n.left);
            queue.add(n.right);
        }
        while (parts.get(parts.size() - 1).equals("null")) {
            parts.remove(parts.size() - 1);
        }
        return "[" + String.join(", ", parts) + "]";
    }

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        br = new BufferedReader(new InputStreamReader(System.in));
        int t = Integer.parseInt(br.readLine().trim());
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            String line = br.readLine();
            TreeNode root = parseTree(line);
            try {
                List<List<Integer>> rows = new Solution().levels(root);
                StringBuilder line2 = new StringBuilder("[");
                for (int i = 0; i < rows.size(); i++) {
                    if (i > 0) {
                        line2.append(", ");
                    }
                    line2.append(rows.get(i).toString());
                }
                out.append(line2).append("]\n");
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
