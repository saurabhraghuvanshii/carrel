import java.io.*;
import java.util.*;

class Node {
    int val;
    List<Node> neighbors = new ArrayList<>();
    Node(int val) { this.val = val; }
}

// Reads T cases from stdin and prints one line per case. Each case is "n m", then m edges "u v".
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

    // Reads m edges of width numbers each: "u v" or "u v w".
    static int[][] readEdges(int m, int width) throws IOException {
        int[][] edges = new int[m][width];
        for (int[] e : edges) {
            for (int j = 0; j < width; j++) {
                e[j] = nextInt();
            }
        }
        return edges;
    }

    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        br = new BufferedReader(new InputStreamReader(System.in));
        int t = nextInt();
        StringBuilder out = new StringBuilder();
        for (int c = 0; c < t; c++) {
            int n = nextInt();
            int m = nextInt();
            int[][] edges = readEdges(m, 2);
            Node[] nodes = new Node[n];
            for (int i = 0; i < n; i++) {
                nodes[i] = new Node(i);
            }
            for (int[] e : edges) {
                nodes[e[0]].neighbors.add(nodes[e[1]]);
                nodes[e[1]].neighbors.add(nodes[e[0]]);
            }
            try {
                Node copy = new Solution().cloneGraph(nodes[0]);
                Set<Node> originals = Collections.newSetFromMap(new IdentityHashMap<>());
                originals.addAll(Arrays.asList(nodes));
                Node[] byVal = new Node[n];
                Set<Node> seen = Collections.newSetFromMap(new IdentityHashMap<>());
                ArrayDeque<Node> queue = new ArrayDeque<>();
                String problem = copy == null ? "nothing was returned" : null;
                if (copy != null) {
                    queue.add(copy);
                    seen.add(copy);
                }
                while (problem == null && !queue.isEmpty()) {
                    Node x = queue.poll();
                    if (originals.contains(x)) {
                        problem = "the copy uses original node " + x.val;
                    } else if (x.val < 0 || x.val >= n || byVal[x.val] != null) {
                        problem = "node " + x.val + " is copied more than once or has a wrong value";
                    } else {
                        byVal[x.val] = x;
                        for (Node y : x.neighbors) {
                            if (y != null && seen.add(y)) {
                                queue.add(y);
                            }
                        }
                    }
                }
                for (int i = 0; problem == null && i < n; i++) {
                    if (byVal[i] == null) {
                        problem = "node " + i + " is missing from the copy";
                    }
                }
                if (problem != null) {
                    out.append("not a copy: ").append(problem).append('\n');
                } else {
                    StringBuilder line = new StringBuilder("[");
                    for (int i = 0; i < n; i++) {
                        if (i > 0) {
                            line.append(", ");
                        }
                        line.append('[');
                        for (int j = 0; j < byVal[i].neighbors.size(); j++) {
                            if (j > 0) {
                                line.append(", ");
                            }
                            line.append(byVal[i].neighbors.get(j).val);
                        }
                        line.append(']');
                    }
                    out.append(line).append("]\n");
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
