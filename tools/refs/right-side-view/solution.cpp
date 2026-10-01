#include <vector>
using namespace std;

// Depth first, right child first: the first node seen at each depth is visible.
static void walk(TreeNode* n, size_t depth, vector<int>& view) {
    if (!n) return;
    if (depth == view.size()) view.push_back(n->val);
    walk(n->right, depth + 1, view);
    walk(n->left, depth + 1, view);
}

vector<int> rightView(TreeNode* root) {
    vector<int> view;
    walk(root, 0, view);
    return view;
}
