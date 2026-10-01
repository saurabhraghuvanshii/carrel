# Driver snippets

Copy these into a pack's `driver.java` or `driver.cpp`. They are not included automatically. The formats are described in the top-level `README.md`.

Every driver:

- reads `T`, then `T` cases, and prints exactly one line per case,
- reads all of a case's input before calling the learner's code, so an exception cannot leave the input half read,
- prints `ERROR <message>` for a case that throws,
- never prints an empty line (the runner treats missing lines at the end as a crash). Lists and arrays print as `[a, b, c]`, empty as `[]`,
- sends the learner's own prints to stderr.

The working examples are the packs `pair-with-target-sum` (array), `reverse-list` (linked list), `max-depth` (tree) and `lru-cache` (call sequence).

## Java

Token reader, used by every format except trees:

```java
static BufferedReader br;
static StringTokenizer st;

static String next() throws IOException {
    while (st == null || !st.hasMoreTokens()) {
        st = new StringTokenizer(br.readLine());
    }
    return st.nextToken();
}

static int nextInt() throws IOException { return Integer.parseInt(next()); }
```

Array (`n`, then n values; an empty array is `0` and an empty line):

```java
int n = nextInt();
int[] a = new int[n];
for (int i = 0; i < n; i++) a[i] = nextInt();
```

Linked list (same as an array; for cycle problems a second line with the index the tail points to, or `-1`):

```java
class ListNode {
    int val;
    ListNode next;
    ListNode(int val) { this.val = val; }
}

int n = nextInt();
ListNode[] nodes = new ListNode[n];
for (int i = 0; i < n; i++) nodes[i] = new ListNode(nextInt());
for (int i = 0; i + 1 < n; i++) nodes[i].next = nodes[i + 1];
ListNode head = n == 0 ? null : nodes[0];
// cycle problems only:
int pos = nextInt();
if (pos >= 0) nodes[n - 1].next = nodes[pos];
```

Printing a list, with a guard against a loop in the returned list:

```java
StringBuilder line = new StringBuilder("[");
int count = 0;
for (ListNode p = result; p != null; p = p.next) {
    if (++count > n) throw new IllegalStateException("the returned list has more nodes than the input. Is there a loop?");
    if (count > 1) line.append(", ");
    line.append(p.val);
}
out.append(line).append("]\n");
```

Binary tree (one line in level order, `null` for a missing child, `null` alone for an empty tree). Read the case with `br.readLine()`, and read `T` with `br.readLine()` too, not the token reader:

```java
class TreeNode {
    int val;
    TreeNode left, right;
    TreeNode(int val) { this.val = val; }
}

static TreeNode parse(String line) {
    String[] tok = line.trim().split("\\s+");
    if (tok[0].isEmpty() || tok[0].equals("null")) return null;
    TreeNode root = new TreeNode(Integer.parseInt(tok[0]));
    ArrayDeque<TreeNode> queue = new ArrayDeque<>();
    queue.add(root);
    for (int i = 1; i < tok.length; ) {
        TreeNode node = queue.poll();
        if (i < tok.length && !tok[i].equals("null")) { node.left = new TreeNode(Integer.parseInt(tok[i])); queue.add(node.left); }
        i++;
        if (i < tok.length && !tok[i].equals("null")) { node.right = new TreeNode(Integer.parseInt(tok[i])); queue.add(node.right); }
        i++;
    }
    return root;
}
```

Graph (`n m`, then m lines `u v` or `u v w`):

```java
int n = nextInt(), m = nextInt();
int[][] edges = new int[m][];
for (int i = 0; i < m; i++) edges[i] = new int[] {nextInt(), nextInt()}; // add nextInt() for w
```

Grid (`rows cols`, then the rows):

```java
int rows = nextInt(), cols = nextInt();
int[][] grid = new int[rows][cols];
for (int r = 0; r < rows; r++) for (int c = 0; c < cols; c++) grid[r][c] = nextInt();
```

Call sequence (`k`, then k lines `name arg ...`; read every call first, then run them; output one result per call joined by spaces, `null` for calls that return nothing). See `lru-cache/driver.java` for the full loop.

## C++

Every C++ driver includes the standard headers it needs, then `using namespace std;`, then the node types, then `#include "solution.cpp"`.

Array and linked list:

```cpp
struct ListNode {
    int val;
    ListNode* next;
    ListNode(int v) : val(v), next(nullptr) {}
};

int n;
cin >> n;
vector<ListNode*> nodes(n);
for (auto& p : nodes) { int v; cin >> v; p = new ListNode(v); }
for (int i = 0; i + 1 < n; i++) nodes[i]->next = nodes[i + 1];
ListNode* head = n == 0 ? nullptr : nodes[0];
// cycle problems only:
int pos;
cin >> pos;
if (pos >= 0) nodes[n - 1]->next = nodes[pos];
```

Binary tree. Read `T` with `getline` too, then one `getline` per case:

```cpp
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode(int v) : val(v), left(nullptr), right(nullptr) {}
};

TreeNode* parseTree(const string& line) {
    istringstream in(line);
    vector<string> tok;
    for (string s; in >> s;) tok.push_back(s);
    if (tok.empty() || tok[0] == "null") return nullptr;
    TreeNode* root = new TreeNode(stoi(tok[0]));
    queue<TreeNode*> q;
    q.push(root);
    for (size_t i = 1; i < tok.size();) {
        TreeNode* node = q.front();
        q.pop();
        if (i < tok.size() && tok[i] != "null") { node->left = new TreeNode(stoi(tok[i])); q.push(node->left); }
        i++;
        if (i < tok.size() && tok[i] != "null") { node->right = new TreeNode(stoi(tok[i])); q.push(node->right); }
        i++;
    }
    return root;
}
```

Graph and grid read the same way as in Java, with `cin >>`.

Call sequence: read all names and arguments into vectors, then run them inside one `try`, writing results into an `ostringstream` so a throw replaces the whole line with `ERROR ...`. See `lru-cache/driver.cpp`.
