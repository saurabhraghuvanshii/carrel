import java.util.*;

// Level order with "#" for a missing child, read back with a queue.
class Codec {
    public String serialize(TreeNode root) {
        StringBuilder b = new StringBuilder();
        LinkedList<TreeNode> queue = new LinkedList<>();
        queue.add(root);
        while (!queue.isEmpty()) {
            TreeNode n = queue.poll();
            if (n == null) {
                b.append("# ");
            } else {
                b.append(n.val).append(' ');
                queue.add(n.left);
                queue.add(n.right);
            }
        }
        return b.toString();
    }

    public TreeNode deserialize(String text) {
        String[] tok = text.trim().split(" ");
        if (tok[0].equals("#")) {
            return null;
        }
        TreeNode root = new TreeNode(Integer.parseInt(tok[0]));
        ArrayDeque<TreeNode> queue = new ArrayDeque<>();
        queue.add(root);
        for (int i = 1; i < tok.length; i += 2) {
            TreeNode n = queue.poll();
            if (!tok[i].equals("#")) {
                n.left = new TreeNode(Integer.parseInt(tok[i]));
                queue.add(n.left);
            }
            if (!tok[i + 1].equals("#")) {
                n.right = new TreeNode(Integer.parseInt(tok[i + 1]));
                queue.add(n.right);
            }
        }
        return root;
    }
}
