import java.io.*;
import java.util.*;

class TreeNode {
    int val;
    TreeNode left, right;
    TreeNode(int val) { this.val = val; }
}

// Flushes after every case so a crash cannot hide finished results.
// Reads T cases from stdin, each one tree in level order on one line with
// "null" for a missing child. Prints one line per case.
// Anything the learner prints goes to stderr so it cannot break the results.
public class Main {
    static TreeNode parse(String line) {
        String[] tok = line.trim().split("\\s+");
        if (tok.length == 0 || tok[0].isEmpty() || tok[0].equals("null")) {
            return null;
        }
        TreeNode root = new TreeNode(Integer.parseInt(tok[0]));
        ArrayDeque<TreeNode> queue = new ArrayDeque<>();
        queue.add(root);
        int i = 1;
        while (i < tok.length) {
            TreeNode node = queue.poll();
            if (i < tok.length && !tok[i].equals("null")) {
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

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        int t = Integer.parseInt(br.readLine().trim());
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            TreeNode root = parse(br.readLine());
            try {
                out.append(new Solution().maxDepth(root)).append('\n');
            } catch (Throwable e) {
                out.append("ERROR ").append(e).append('\n');
            }
            real.print(out);
            real.flush();
            out.setLength(0);
        }
    }
}
