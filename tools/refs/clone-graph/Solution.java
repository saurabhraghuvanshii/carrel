import java.util.*;

class Solution {
    public Node cloneGraph(Node node) {
        if (node == null) {
            return null;
        }
        Map<Node, Node> copies = new HashMap<>();
        copies.put(node, new Node(node.val));
        ArrayDeque<Node> queue = new ArrayDeque<>();
        queue.add(node);
        while (!queue.isEmpty()) {
            Node x = queue.poll();
            for (Node y : x.neighbors) {
                if (!copies.containsKey(y)) {
                    copies.put(y, new Node(y.val));
                    queue.add(y);
                }
                copies.get(x).neighbors.add(copies.get(y));
            }
        }
        return copies.get(node);
    }
}
