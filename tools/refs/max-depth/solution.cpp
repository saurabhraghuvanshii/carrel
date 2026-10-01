#include <queue>
using namespace std;

int maxDepth(TreeNode* root) {
    int depth = 0;
    queue<TreeNode*> level;
    if (root != nullptr) level.push(root);
    while (!level.empty()) {
        depth++;
        for (int i = (int)level.size(); i > 0; i--) {
            TreeNode* n = level.front();
            level.pop();
            if (n->left != nullptr) level.push(n->left);
            if (n->right != nullptr) level.push(n->right);
        }
    }
    return depth;
}
