import java.util.*;

// Ported from the owner's Graphs/BreadthFirstSearch.java (hasPath), unchanged; the graph is
// built from the input edges instead of createGraph.
class Solution {
    static class Edge {
        int src;
        int dest;
        int wt;

        Edge(int s, int d, int w) {
            this.src = s;
            this.dest = d;
            this.wt = w;
        }
    }

    @SuppressWarnings("unchecked")
    private static ArrayList<Edge>[] build(int n, int[][] edges) {
        ArrayList<Edge>[] graph = new ArrayList[n];
        for (int i = 0; i < n; i++) {
            graph[i] = new ArrayList<>();
        }
        for (int[] e : edges) {
            graph[e[0]].add(new Edge(e[0], e[1], 1));
            graph[e[1]].add(new Edge(e[1], e[0], 1));
        }
        return graph;
    }

    public boolean pathExists(int n, int[][] edges, int source, int target) {
        return hasPath(build(n, edges), source, target, new boolean[n]);
    }

    private boolean hasPath(ArrayList<Edge>[] graph, int src, int dest, boolean[] vis) {
        if (src == dest) {
            return true;
        }
        vis[src] = true;
        for (int i = 0; i < graph[src].size(); i++) {
            Edge e = graph[src].get(i);
            if (!vis[e.dest] && hasPath(graph, e.dest, dest, vis)) {
                return true;
            }
        }
        return false;
    }
}
