#include <queue>
#include <unordered_map>
#include <vector>
using namespace std;

Node* cloneGraph(Node* node) {
    if (!node) return nullptr;
    unordered_map<Node*, Node*> copies = {{node, new Node(node->val)}};
    queue<Node*> q;
    q.push(node);
    while (!q.empty()) {
        Node* x = q.front();
        q.pop();
        for (Node* y : x->neighbors) {
            if (!copies.count(y)) {
                copies[y] = new Node(y->val);
                q.push(y);
            }
            copies[x]->neighbors.push_back(copies[y]);
        }
    }
    return copies[node];
}
